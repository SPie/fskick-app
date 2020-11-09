<?php

namespace Tests\Helper;

use App\Console\TableHelperFactory;
use Illuminate\Console\Command;
use Illuminate\Console\OutputStyle;
use Mockery as m;
use Mockery\MockInterface;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

/**
 * Trait ConsoleHelper
 *
 * @package Tests\Helper
 */
trait ConsoleHelper
{
    /**
     * @param Command        $command
     * @param InputInterface $input
     * @param OutputStyle    $output
     *
     * @return Command
     */
    private function setInputAndOutput(Command $command, InputInterface $input, OutputStyle $output): Command
    {
        $command->setInput($input);
        $command->setOutput($output);

        return $command;
    }

    /**
     * @return InputInterface|MockInterface
     */
    private function createInput(): InputInterface
    {
        return m::spy(InputInterface::class);
    }

    /**
     * @param InputInterface|MockInterface $input
     * @param string|null                  $argument
     * @param string|null                  $key
     *
     * @return $this
     */
    private function mockInputGetArgument(MockInterface $input, ?string $argument, ?string $key): self
    {
        $input
            ->shouldReceive('getArgument')
            ->with($key)
            ->andReturn($argument);

        return $this;
    }

    /**
     * @return OutputStyle|MockInterface
     */
    private function createOutputStyle(): OutputStyle
    {
        return m::spy(OutputStyle::class);
    }

    /**
     * @param OutputStyle|MockInterface $outputStyle
     * @param string                    $line
     *
     * @param int|null                  $verbosity
     *
     * @return $this
     */
    private function assertOutputStyleWriteln(MockInterface $outputStyle, string $line, int $verbosity): self
    {
        $outputStyle
            ->shouldHaveReceived('writeln')
            ->with($line, $verbosity)
            ->once();

        return $this;
    }

    /**
     * @return Table|MockInterface
     */
    private function createTableHelper(): Table
    {
        return m::spy(Table::class);
    }

    /**
     * @param Table|MockInterface $tableHelper
     *
     * @return $this
     */
    private function assertTableHelperRender(MockInterface $tableHelper): self
    {
        $tableHelper->shouldHaveReceived('render')->once();

        return $this;
    }

    /**
     * @return TableHelperFactory|MockInterface
     */
    private function createTableHelperFactory(): TableHelperFactory
    {
        return m::spy(TableHelperFactory::class);
    }

    /**
     * @param MockInterface   $tableHelperFactory
     * @param Table           $tableHelper
     * @param OutputInterface $output
     * @param array           $headers
     * @param array           $rows
     *
     * @return $this
     */
    private function mockTableRendererRenderCreate(
        MockInterface $tableHelperFactory,
        Table $tableHelper,
        OutputInterface $output,
        array $headers,
        array $rows
    ): self {
        $tableHelperFactory
            ->shouldReceive('create')
            ->with($output, $headers, $rows)
            ->andReturn($tableHelper);

        return $this;
    }
}
