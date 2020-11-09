<?php

namespace Tests\Unit\Console;

use App\Console\TableHelperFactory;
use Symfony\Component\Console\Exception\InvalidArgumentException;
use Symfony\Component\Console\Helper\Table;
use Tests\Helper\ConsoleHelper;
use Tests\TestCase;

/**
 * Class TableHelperFactoryTest
 *
 * @package Tests\Unit\Console
 */
final class TableHelperFactoryTest extends TestCase
{
    use ConsoleHelper;

    //region Tests

    /**
     * @return void
     */
    public function testCreate(): void
    {
        $output = $this->createOutputStyle();
        $headers = [$this->getFaker()->word];
        $rows = [[$this->getFaker()->word]];

        $this->assertEquals(
            (new Table($output))
                ->setHeaders($headers)
                ->setRows($rows),
            $this->getTableHelperFactory()->create($output, $headers, $rows)
        );
    }

    /**
     * @return void
     */
    public function testCreateWithoutArrayAsRow(): void
    {
        $this->expectException(InvalidArgumentException::class);

        $this->getTableHelperFactory()->create($this->createOutputStyle(), [$this->getFaker()->word], [$this->getFaker()->word]);
    }

    //endregion

    /**
     * @return TableHelperFactory
     */
    private function getTableHelperFactory(): TableHelperFactory
    {
        return new TableHelperFactory();
    }
}
