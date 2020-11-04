<?php

namespace Tests\Unit\Models;

use App\Models\DoctrineDatabaseHandler;
use App\Models\Model;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\ORM\EntityManager;
use Doctrine\ORM\Persisters\Entity\EntityPersister;
use Doctrine\ORM\UnitOfWork;
use Mockery as m;
use Mockery\MockInterface;
use Tests\Helper\ModelHelper;
use Tests\TestCase;

/**
 * Class DoctrineDatabaseHandlerTest
 *
 * @package Tests\Unit\Models
 */
final class DoctrineDatabaseHandlerTest extends TestCase
{
    use ModelHelper;

    //region Tests

    /**
     * @param bool $withModel
     *
     * @return array
     */
    private function setUpFindTest(bool $withModel = true): array
    {
        $id = $this->getFaker()->numberBetween();
        $className = $this->getFaker()->word;
        $model = $this->createModel();
        $entityManager = $this->createEntityManager();
        $this->mockEntityManagerFind($entityManager, $withModel ? $model : null, $className, $id);
        $doctrineDatabaseHandler = $this->getDoctrineDatabaseHandler($entityManager, $className);

        return [$doctrineDatabaseHandler, $id, $model];
    }

    /**
     * @return void
     */
    public function testFind(): void
    {
        [$doctrineDatabaseHandler, $id, $model] = $this->setUpFindTest();

        $this->assertEquals($model, $doctrineDatabaseHandler->find($id));
    }

    /**
     * @return void
     */
    public function testFindWithoutModel(): void
    {
        [$doctrineDatabaseHandler, $id] = $this->setUpFindTest(false);

        $this->assertNull($doctrineDatabaseHandler->find($id));
    }

    private function setUpLoadTest(bool $withModel = true): array
    {
        $criteria = [$this->getFaker()->word => $this->getFaker()->word];
        $className = $this->getFaker()->word;
        $model = $this->createModel();
        $entityPersister = $this->createEntityPersister();
        $this->mockEntityPersisterLoad($entityPersister, $withModel ? $model : null, $criteria);
        $databaseHandler = $this->getDoctrineDatabaseHandler(
            $this->createEntityManagerWithEntityPersister($className, $entityPersister),
            $className
        );

        return [$databaseHandler, $criteria, $model];
    }

    /**
     * @return void
     */
    public function testLoad(): void
    {
        [$databaseHandler, $criteria, $model] = $this->setUpLoadTest();

        $this->assertEquals($model, $databaseHandler->load($criteria));
    }

    /**
     * @return void
     */
    public function testLoadWithoutModel(): void
    {
        [$databaseHandler, $criteria] = $this->setUpLoadTest(false);

        $this->assertEmpty($databaseHandler->load($criteria));
    }

    private function setUpLoadAllTest(bool $withModels = true, bool $withoutOptionalParameters = false): array
    {
        $criteria = [$this->getFaker()->word => $this->getFaker()->word];
        $orderBy = $withoutOptionalParameters ? [] : [$this->getFaker()->word => $this->getFaker()->word];
        $limit = $withoutOptionalParameters ? null : $this->getFaker()->numberBetween();
        $offset = $withoutOptionalParameters ? null : $this->getFaker()->numberBetween();
        $models = [$this->createModel()];
        $entityPersister = $this->createEntityPersister();
        $this->mockEntityPersisterLoadAll($entityPersister, $withModels ? $models : [], $criteria, $orderBy, $limit, $offset);
        $className = $this->getFaker()->word;
        $entityManager = $this->createEntityManagerWithEntityPersister($className, $entityPersister);
        $databaseHandler = $this->getDoctrineDatabaseHandler($entityManager, $className);

        return [$databaseHandler, $criteria, $orderBy, $limit, $offset, $models];
    }

    /**
     * @return void
     */
    public function testLoadAll(): void
    {
        [$databaseHandler, $criteria, $orderBy, $limit, $offset, $models] = $this->setUpLoadAllTest();

        $this->assertEquals(new ArrayCollection($models), $databaseHandler->loadAll($criteria, $orderBy, $limit, $offset));
    }

    /**
     * @return void
     */
    public function testLoadAllWithoutModels(): void
    {
        [$databaseHandler, $criteria, $orderBy, $limit, $offset] = $this->setUpLoadAllTest(false);

        $this->assertEquals(new ArrayCollection(), $databaseHandler->loadAll($criteria, $orderBy, $limit, $offset));
    }

    /**
     * @return void
     */
    public function testLoadAllWithoutOptionalParameters(): void
    {
        [$databaseHandler, $criteria, $orderBy, $limit, $offset, $models] = $this->setUpLoadAllTest(true, true);

        $this->assertEquals(new ArrayCollection($models), $databaseHandler->loadAll($criteria));
    }

    private function setUpSaveTest(): array
    {
        $model = $this->createModel();
        $entityManager = $this->createEntityManager();
        $databaseHandler = $this->getDoctrineDatabaseHandler($entityManager);

        return [$databaseHandler, $model, $entityManager];
    }

    /**
     * @return void
     */
    public function testSave(): void
    {
        [$databaseHandler, $model, $entityManager] = $this->setUpSaveTest();

        $this->assertEquals($model, $databaseHandler->save($model));
        $this
            ->assertEntityManagerPersist($entityManager, $model)
            ->assertEntityManagerFlush($entityManager);
    }

    /**
     * @return void
     */
    public function testSaveWithoutFlush(): void
    {
        [$databaseHandler, $model, $entityManager] = $this->setUpSaveTest();

        $databaseHandler->save($model, false);

        $entityManager->shouldNotHaveReceived('flush');
    }

    private function setUpDeleteTest(): array
    {
        $model = $this->createModel();
        $entityManager = $this->createEntityManager();
        $databaseHandler = $this->getDoctrineDatabaseHandler($entityManager);

        return [$databaseHandler, $model, $entityManager];
    }

    /**
     * @return void
     */
    public function testDelete(): void
    {
        [$databaseHandler, $model, $entityManager] = $this->setUpDeleteTest();

        $databaseHandler->delete($model);

        $this
            ->assertEntityManagerRemove($entityManager, $model)
            ->assertEntityManagerFlush($entityManager);
    }

    /**
     * @return void
     */
    public function testDeleteWithoutFlush(): void
    {
        [$databaseHandler, $model, $entityManager] = $this->setUpDeleteTest();

        $databaseHandler->delete($model, false);

        $entityManager->shouldNotHaveReceived('flush');
    }

    /**
     * @return void
     */
    public function testFlush(): void
    {
        $entityManager = $this->createEntityManager();

        $this->getDoctrineDatabaseHandler($entityManager)->flush();

        $this->assertEntityManagerFlush($entityManager);
    }

    //endregion

    /**
     * @param EntityManager|null $entityManager
     * @param string|null        $className
     *
     * @return DoctrineDatabaseHandler
     */
    private function getDoctrineDatabaseHandler(
        EntityManager $entityManager = null,
        string $className = null
    ): DoctrineDatabaseHandler {
        return new DoctrineDatabaseHandler(
            $entityManager ?: $this->createEntityManager(),
            $className ?: $this->getFaker()->word
        );
    }

    /**
     * @return EntityManager|MockInterface
     */
    private function createEntityManager(): EntityManager
    {
        return m::spy(EntityManager::class);
    }

    /**
     * @param string|null          $className
     * @param EntityPersister|null $entityPersister
     *
     * @return EntityManager|MockInterface
     */
    private function createEntityManagerWithEntityPersister(
        string $className = null,
        EntityPersister $entityPersister = null
    ): EntityManager {
        $unitOfWork = $this->createUnitOfWork();
        $this->mockUnitOfWorkGetEntityPersister(
            $unitOfWork,
            $entityPersister ?: $this->createEntityPersister(),
            $className ?: $this->getFaker()->word
        );

        $entityManager = $this->createEntityManager();
        $this->mockEntityManagerGetUniOfWork($entityManager, $unitOfWork);

        return $entityManager;
    }

    /**
     * @param EntityManager|MockInterface $entityManager
     * @param UnitOfWork                  $unitOfWork
     *
     * @return $this
     */
    private function mockEntityManagerGetUniOfWork(MockInterface $entityManager, UnitOfWork $unitOfWork): self
    {
        $entityManager
            ->shouldReceive('getUnitOfWork')
            ->andReturn($unitOfWork);

        return $this;
    }

    /**
     * @param EntityManager|MockInterface $entityManager
     * @param Model|null                  $model
     * @param string                      $className
     * @param int                         $id
     *
     * @return $this
     */
    private function mockEntityManagerFind(MockInterface $entityManager, ?Model $model, string $className, int $id): self
    {
        $entityManager
            ->shouldReceive('find')
            ->with($className, $id)
            ->andReturn($model);

        return $this;
    }

    /**
     * @return EntityPersister|MockInterface
     */
    private function createEntityPersister(): EntityPersister
    {
        return m::spy(EntityPersister::class);
    }

    /**
     * @param EntityPersister|MockInterface $entityPersister
     * @param Model|null                    $model
     * @param array                         $criteria
     *
     * @return $this
     */
    private function mockEntityPersisterLoad(MockInterface $entityPersister, ?Model $model, array $criteria): self
    {
        $entityPersister
            ->shouldReceive('load')
            ->with($criteria)
            ->andReturn($model);

        return $this;
    }

    private function mockEntityPersisterLoadAll(
        MockInterface $entityPersister,
        array $models,
        array $criteria,
        array $orderBy,
        ?int $limit,
        ?int $offset
    ): self {
        $entityPersister
            ->shouldReceive('loadAll')
            ->with($criteria, $orderBy, $limit, $offset)
            ->andReturn($models);

        return $this;
    }

    /**
     * @param EntityManager|MockInterface $entityManager
     * @param Model                       $model
     *
     * @return $this
     */
    private function assertEntityManagerPersist(MockInterface $entityManager, Model $model): self
    {
        $entityManager
            ->shouldHaveReceived('persist')
            ->with($model)
            ->once();

        return $this;
    }

    /**
     * @param EntityManager|MockInterface $entityManager
     * @param Model                       $model
     *
     * @return $this
     */
    private function assertEntityManagerRemove(MockInterface $entityManager, Model $model): self
    {
        $entityManager
            ->shouldHaveReceived('remove')
            ->with($model)
            ->once();

        return $this;
    }

    /**
     * @param EntityManager|MockInterface $entityManager
     *
     * @return $this
     */
    private function assertEntityManagerFlush(MockInterface $entityManager): self
    {
        $entityManager
            ->shouldHaveReceived('flush')
            ->once();

        return $this;
    }

    /**
     * @return UnitOfWork|MockInterface
     */
    private function createUnitOfWork(): UnitOfWork
    {
        return m::spy(UnitOfWork::class);
    }

    /**
     * @param UnitOfWork|MockInterface $unitOfWork
     * @param EntityPersister          $entityPersister
     * @param string                   $className
     *
     * @return $this
     */
    private function mockUnitOfWorkGetEntityPersister(
        MockInterface $unitOfWork,
        EntityPersister $entityPersister,
        string $className
    ): self {
        $unitOfWork
            ->shouldReceive('getEntityPersister')
            ->with($className)
            ->andReturn($entityPersister);

        return $this;
    }
}
